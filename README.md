# HashRing implements consistent hashing with partitioning

# Flow Diagram:

# +-----------+        +---------------+        +----------------+
# | Client    +------->+ HashRing      +------->+ Node Responsible|
# | PUT/GET   |        | hashKey(key)  |        | for that hash   |
# +-----------+        +---------------+        +----------------+
#                                                |
#                                                v
#                                       +------------------+
#                                       | Store/Retrieve KV|
#                                       +------------------+

# Description:
# 1. The client issues a PUT or GET request with a key.
# 2. HashRing hashes the key using SHA-1 (first 4 bytes).
# 3. Sorted node hashes determine which node the key maps to.
# 4. Data is routed to the appropriate node.
# 5. Partitioning is achieved by consistent hashing distributing keys across nodes.