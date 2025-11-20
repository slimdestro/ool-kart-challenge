
# **OOL Coupon Validation and API Suite**

```
                                                           ┌────────────────────────────┐
                                                           │      .gz Coupon Databases   │
                                                           └───────────────┬────────────┘
                                                                           │
                                                             Multi-threaded Loader
                                                                           │
                                                           ┌───────────────▼──────────────┐
                                                           │     Hashing: FNV-1a (uint32)  │
                                                           └───────────────┬──────────────┘
                                                                           │
                                                               In-Memory Hash Index (map)
                                                                           │
                                             ┌─────────────────────────────┼─────────────────────────────┐
                                             │                             │                             │
                                     LRU Cache (Fast Path)       Validation Engine              File-Mask Logic
                                             │                             │                             │
                                             └─────────────────────────────┼─────────────────────────────┘
                                                                           │
                                                                     `/api/order`
```

---

# 📸 **Screenshots**

### **System Running**

![Screenshot 1](assets/Screenshot%202025-11-20%20234535.png)

### **Validation Logs**

![Screenshot 2](assets/Screenshot%202025-11-20%20234834.png)

### **Performance / Execution Output**

![Screenshot 3](assets/Screenshot%202025-11-20%20234903.png)

---

# ⚙️ **Quick Start**

```sh
ool-server.exe \
  -f1 coupons/couponbase1.gz \
  -f2 coupons/couponbase2.gz \
  -f3 coupons/couponbase3.gz \
  -port 8080
```

---

# 📡 **API Examples**

| Endpoint     | Method | Description                    | Example                          |
| ------------ | ------ | ------------------------------ | -------------------------------- |
| `/api/order` | POST   | Validate a coupon code         | `{ "coupon": "SAVE100" }`        |
| `/api/ping`  | GET    | Health check                   | `"pong"`                         |
| `/api/stats` | GET    | Memory, cache hit-rate, uptime | `{ "hits": 1234, "misses": 98 }` |

---

# 🚀 **Core Features**

## **1. Memory-Optimized Coupon Validation**

### Problem

Storing millions of codes as `map[string]byte` consumed **6–7 GB RAM**.

### Solution

✔ Convert coupon strings to **FNV-1a uint32 hashes**
✔ Store in a compact **`map[uint32]byte`** index
✔ Dramatically lower memory usage

---

## **2. Concurrent Indexing & High-Speed Loading**

* Multi-threaded loading for `.gz` files
* Significant reduction in startup time
* Streaming I/O + `WaitGroup` pipeline

---

## **3. Multi-Source Validation Logic**

A coupon is considered valid only if:

1. The hash exists in the in-memory index **and**
2. It appears in **2 or more** of the loaded databases

Prevents single-source false positives and increases trust.

---

## **4. LRU Cache Acceleration**

* Hot coupon lookups skip hashing + DB checks
* Sub-millisecond validation for cached entries
* Fully thread-safe

---

# 📊 **Performance**

| Metric                                  | Result                            |
| --------------------------------------- | --------------------------------- |
| Startup load (3× ~3GB compressed files) | **Fast, multi-threaded**          |
| Memory footprint                        | **~600 MB**                       |
| Throughput                              | **1M+ requests/min** (local test) |
| Latency (cached)                        | **< 1 ms**                        |
