package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"

    pb "github.com/BeratHundurel/order-api/currency-api/proto"

    "google.golang.org/grpc"
)

const fixerAPIURL = "http://data.fixer.io/api/latest"

type server struct {
    pb.UnimplementedCurrencyConverterServer
}

func (s *server) ConvertCurrency(ctx context.Context, req *pb.ConvertCurrencyRequest) (*pb.ConvertCurrencyResponse, error) {
    conversionRate, err := getConversionRate(req.FromCurrency, req.ToCurrency)
    if err != nil {
        return nil, err
    }
    convertedPrice := req.Price * conversionRate

    return &pb.ConvertCurrencyResponse{
        ConvertedPrice: convertedPrice,
        ToCurrency:     req.ToCurrency,
    }, nil
}

func getConversionRate(fromCurrency, toCurrency string) (float32, error) {
    apiKey := os.Getenv("FIXER_API_KEY")
    url := fmt.Sprintf("%s?access_key=%s&base=%s&symbols=%s", fixerAPIURL, apiKey, fromCurrency, toCurrency)

    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch data from Fixer API: %v", err)
    }
    defer resp.Body.Close()

    var result struct {
        Success    bool                `json:"success"`
        Timestamp  int64               `json:"timestamp"`
        Base       string              `json:"base"`
        Date       string              `json:"date"`
        Rates      map[string]float32  `json:"rates"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, fmt.Errorf("failed to decode Fixer API response: %v", err)
    }

    if !result.Success {
        return 0, fmt.Errorf("failed to fetch conversion rates: API returned an unsuccessful response")
    }

    rate, ok := result.Rates[toCurrency]
    if !ok {
        return 0, fmt.Errorf("conversion rate for %s not found", toCurrency)
    }

    return rate, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterCurrencyConverterServer(s, &server{})
    log.Printf("Server is running on port :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

