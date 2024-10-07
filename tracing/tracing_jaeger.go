package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initJeagerTracer(endpoint string) (func(), error) {
	log.Printf("Initializing tracer with Jaeger endpoint: %s", endpoint)

	// 初始化 Jaeger 導出器
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, fmt.Errorf("創建 Jaeger 導出器失敗: %w", err)
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

func TracingJeager() {
	// Jaeger 的端點
	// endpoint := "http://localhost:14268/api/traces"
	endpoint := "http://jaeger:14268/api/traces"
	shutdown, err := initJeagerTracer(endpoint)
	if err != nil {
		log.Fatalf("初始化追踪器失敗: %v", err)
	}
	defer shutdown()

	tracer := otel.Tracer("example-tracer")

	xTimes := 20
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
