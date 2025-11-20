package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"ool/internal/app/coupon"
	"ool/internal/handler"
	"ool/pkg/server"
)

func memStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Current Memory: PermanentAlloc=%.2fMB, SystemMem=%.2fMB",
		float64(m.Alloc)/1024/1024,
		float64(m.Sys)/1024/1024)
}

func main() {
	var f1, f2, f3 string
	var port int

	flag.StringVar(&f1, "f1", "", "Path to couponbase1.gz (Required if using 2 files)")
	flag.StringVar(&f2, "f2", "", "Path to couponbase2.gz (Required)")
	flag.StringVar(&f3, "f3", "", "Path to couponbase3.gz (Optional)")
	flag.IntVar(&port, "port", 8080, "Port")
	flag.Parse()
	var filePaths []string
	if f1 != "" {
		filePaths = append(filePaths, f1)
	}
	if f2 != "" {
		filePaths = append(filePaths, f2)
	}
	if f3 != "" {
		filePaths = append(filePaths, f3)
	}

	numFiles := len(filePaths)
	if numFiles < 2 || numFiles > 3 {
		log.Fatalf("Minimum 2 files required for coupon validation logic")
	}

	//log.Printf("Starting indexing from %d files...", numFiles)
	ci, totalIndexDuration := coupon.BuildIndex(filePaths)

	if ci.Size() == 0 {
		log.Printf("Warning: Indexing complete, but 0 unique coupons were found.")
	} else {
		log.Printf("Indexed %d unique coupons in %s.", ci.Size(), totalIndexDuration.String())
	}

	memStats()

	cache := coupon.NewLRU(100000)
	productRepo := handler.GetProductRepo()
	router := server.NewRouter(ci, cache, productRepo)
	srv := server.NewServer(port, router)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server listening on %d", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown: %v", err)
	}
}
