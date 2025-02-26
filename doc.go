/*
Package cpus supports working with CPU lists and sets, as well as querying the
CPUs currently assigned to processes and tasks and pinning tasks and processes
to specific CPU sets.

Logically, [List] and [Set] are equivalent, as they both represent sets of one
or more logical CPUs. Each logical CPU is identified by their 0-based CPU
number. The difference between List and Set lies in their internal
representations, mirroring different representation forms in the Linux syscalls
and procfs pseudo files.

  - [List] internally stores CPU numbers as ranges, such as 1-4, 8-15.
  - [Set] internally stores CPU numbers as bits in a bytestream, such as (hex)
    ff1e.

[List.Set] converts a List into its corresponding Set. In the opposite
direction, [Set.List] converts a Set into its equivalent List.
*/
package cpus
