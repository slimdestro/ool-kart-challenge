# 💡 Technical Overview: High-Performance Coupon Validation

The **OOL API** is engineered for **high-volume, low-latency coupon validation**, primarily served through the **rate-limited `/api/order`** endpoint.

---


### **1. Coupon validation Memory Optimization (This seems Core Challenge of this test)**

**Problem:**
Directly storing the full coupon dataset as `map[string]byte` resulted in a massive **6–7 GB** memory footprint.

**Solution:**

* Coupons are converted into **`uint32` hashes** using the **FNV-32a** algorithm.
* Validation is performed against an in-memory **`map[uint32]byte` index**, reducing RAM usage dramatically.

---

### **2. Concurrent Indexing & High-Speed Loading**

* The coupon index is built at startup from **2 or more compressed `.gz` sources**.
* **Multi-threaded streaming**, combined with `sync.WaitGroup`, ensures maximum I/O throughput.
* The full index is loaded quickly despite extremely large datasets.

---

### **3. Validation Logic**

A coupon is considered **valid** only when:

1. Its hash exists in the in-memory index, **and**
2. Its associated file-mask indicates it appears in **two or more data sources**.

This prevents false positives and ensures multi-source agreement.

---

### **4. Runtime Acceleration via LRU Cache**

To minimize repetitive hash lookups:

* A high-performance, thread-safe **LRU cache** stores results of recent validations.
* Popular or frequently retried coupon codes are resolved at near-zero cost.
