# Brute approach 

* For calculating mean do not divide at each line, as division operation are costly.
* Execution time    : 90.1746671s


<!-- 
- htop
- time go run main.go
    in real time the time taken to run including overhead of wsl(thats slow due to additional layer of virulisation) 
    // real    22m2.393s
    // user    2m34.303s
    // sys     0m35.381s
-  Measure-Command { go run main.go }

-->
<!--
# Memory mapping  

* mmap is used to avoid frequent io calls. To access any byte an io call is required.
* Virtually map the file to an array. This is done once, and now can acces file like a array
* Execution time    : 8.3929657s  (no processing just counting)
 -->

<!-- 
// real    19m22.253s
// user    0m8.208s
// sys     0m28.267s     -->

# Custom integer parsing + Memory mapping

* Not using float rather taking advantage of values have one decimal place, so directly convert the string to integer(value * 10).   
* Execution time    : 68.7197756s


# Go routines + Memory mapping 

* Divided the file into chunks and process these chunks on differnt cores
* Execution time    : 12.0857479s

# Custom Map with linear probing

* Custom hashmap using FNV-1a hashing algorithm and handling collison through linear probing 
* Execution time    : 4.0971261s
