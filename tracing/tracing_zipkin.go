package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracer(endpoint string) (func(), error) {
	log.Printf("Initializing tracer with Zipkin endpoint: %s", endpoint)

	// 初始化 Zipkin 導出器
	exporter, err := zipkin.New(endpoint)
	if err != nil {
		return nil, fmt.Errorf("創建 Zipkin 導出器失敗: %w", err)
	}

	// 創建資源，描述服務的詳細信息
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("MyService"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("創建資源失敗: %w", err)
	}

	// 創建追踪器提供者
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("關閉追踪器提供者時出錯: %v", err)
		}
	}, nil
}

func doOperationWithCtx(ctx context.Context) error {
	// 模擬操作
	time.Sleep(100 * time.Millisecond)
	return nil
}

func main() {
	// Zipkin 的端點
	endpoint := "http://localhost:9412/api/v2/spans"
	shutdown, err := initTracer(endpoint)
	if err != nil {
		log.Fatalf("初始化追踪器失敗: %v", err)
	}
	defer shutdown()

	tracer := otel.Tracer("example-tracer")

	xTimes := 5
	for i := 0; i < xTimes; i++ {
		// 每次操作一個追踪 span
		ctx, span := tracer.Start(context.Background(), "doOperation")
		err := doOperationWithCtx(ctx)
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}

	fmt.Println("Operations completed")

	// 給導出器一些時間來發送數據
	time.Sleep(2 * time.Second)
}
