# OOL Coupon Validation and Rest of the APIs

                                                       ┌────────────────────────────┐
                                                       │   .gz Coupon Databases      │
                                                       └──────────────┬─────────────┘
                                                                      │
                                                         Multi-threaded Loader
                                                                      │
                                                       ┌──────────────▼──────────────┐
                                                       │  Hashing: FNV-1a (uint32)    │
                                                       └──────────────┬─────────────┘
                                                                      │
                                                         In-Memory Hash Index (map)
                                                                      │
                                                 ┌────────────────────┼────────────────────┐
                                                 │                    │                    │
                                        LRU Cache (Fast Path)    Validation Engine   File-mask Logic
                                                 │                    │                    │
                                                 └────────────────────┼────────────────────┘
                                                                      │
                                                              `/api/order`
```


# 📸 Screenshots

### **System Running**

![Screenshot 1](assets/Screenshot%202025-11-20%20234535.png)

### **Validation Logs**

![Screenshot 2](assets/Screenshot%202025-11-20%20234834.png)

### **Performance / Execution Output**

![Screenshot 3](assets/Screenshot%202025-11-20%20234903.png)

---

# ⚙️ Quick Start

```sh
ool-server.exe \
  -f1 coupons/couponbase1.gz \
  -f2 coupons/couponbase2.gz \
  -f3 coupons/couponbase3.gz \
  -port 8080
```

---

# 📡 **API Examples**

| Endpoint     | Method | Description                                | Example                      |
| ------------ | ------ | ------------------------------------------ | ---------------------------- |
| `/api/order` | `POST` | Validate a coupon code                     | `{ "coupon": "SAVE100" }`    |
| `/api/ping`  | `GET`  | Health check                               | Returns `"pong"`             |
| `/api/stats` | `GET`  | Returns memory, cache hit-rate, and uptime | `{ hits: 1234, misses: 98 }` |

---

# 🚀 Core Features

## **1. Memory-Optimized Coupon Validation**

### **Problem**

Storing millions of codes as `map[string]byte` consumed **6–7 GB RAM**.

### **Solution**

* Convert coupon strings to **FNV-32a `uint32` hashes`**
* Store them in a compact **`map[uint32]byte`** index
* Reduce memory usage **massively**

---

## **2. Concurrent Indexing & High-Speed Loading**

* Load **multiple `.gz` files** in parallel
* Reduce startup time significantly
* Utilize `WaitGroup` + streaming I/O

---

## **3. Multi-Source Validation Logic**

A coupon is only valid if:

1. Its hash exists in the index, **and**
2. It appears in **2+ database files**

This ensures **high confidence** and prevents single-source false positives.

---

## **4. LRU Cache Acceleration**

* Frequent coupon checks skip hashing + lookup
* Near-zero validation time for hot codes
* Thread-safe implementation

---

# 📊 **Performance**

| Metric                                    | Result                            |
| ----------------------------------------- | --------------------------------- |
| Startup load (3× ~3GB compressed sources) | **Fast, multi-threaded**          |
| Memory footprint                          | **Under ~600 MB**                 |
| Throughput                                | **1M+ requests/min** (local test) |
| Latency                                   | **<1 ms** for cached lookups      |
