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
        fmt.Printf("failed to get conversion rate: %v", err)
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
    if apiKey == "" {
        return 0, fmt.Errorf("FIXER_API_KEY environment variable is not set")
    }

    if fromCurrency == "" || toCurrency == "" {
        return 0, fmt.Errorf("fromCurrency or toCurrency is empty")
    }

    url := fmt.Sprintf("%s?access_key=%s&base=%s&symbols=%s", fixerAPIURL, apiKey, fromCurrency, toCurrency)

    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch data from Fixer API: %v", err)
    }
    defer resp.Body.Close()

    var result struct {
        Success    bool                `json:"success"`
        Rates      map[string]float32  `json:"rates"`
        Error      struct {
            Code int    `json:"code"`
            Type string `json:"type"`
        } `json:"error"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, fmt.Errorf("failed to decode Fixer API response: %v", err)
    }

    if !result.Success {
        return 0, fmt.Errorf("failed to fetch conversion rates: %s (code %d)", result.Error.Type, result.Error.Code)
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

