You can use the dd command in Linux to create a binary file with random data. Here's a simple command that generates a 1GB binary file filled with random data:

`dd if=/dev/urandom of=random_data.bin bs=1M count=1024`
