# powerofsingleport
### [TEST] SinglePort vs MultiPort

  Pinning ports to workers is not good idea. When the consumer queue is full, other workers cannot steal busy worker's consumer queue data. In this case the queue has to wait pinned worker which is assigned to itself. Maybe in application level you can attemp to steal other worker's data but that approach will cost you more context-switch. Using one port let the workers to handle scheduling(in this test we let the golang to do that). 
 <br/><br/>__This test shows us the power of singleport against multiport.__


 
![powerofsingleport2](https://github.com/abdullahb53/powerofsingleport/assets/29378922/21842e6d-e276-44a6-9a84-1c76ec0b6a64)

![powerofsingleport](https://github.com/abdullahb53/powerofsingleport/assets/29378922/11262fd1-66eb-4c1e-95f6-cc19e9fdd133)
