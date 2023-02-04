Configure in bridge mode: connection between host and container:
1. Create a veth pair
2. Create a bridge device in the host
3. One side of the veth pair is attached to the bridge and one to the container ns
4. Route traffic from container through veth
