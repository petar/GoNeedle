# What is GoNeedle

GoNeedle is a system for establishing a reliable transport between peers over UDP,
while punching through NATs and firewalls.

GoNeedle is intended to be simple and efficient. Thus, unlike libraries like 
libjingle, GoNeedle uses only one method for punching NATs. This is the most 
general method used in paractice: Two peers (using the help of a server for
coordination) send UDP packets at each other until the firewall is punched.
After successfull punching GoNeedle will (not implemented yet) establish
reliable transport over UDP while using the DCCP protocol for congestion control
and a simple erasure-check scheme for reliability (i.e. for detecting lost packets).

GoNeedle is currently in development/exeprimental stage. Do not install it if 
are looking for something that works out of the box!

# Status

We have implemented the server and the punching procedure and we have made
various tests confirming that it works well. The next stage is to implement
the reliable transport over UDP layer. This is contingent on the completion
of a separate project [GoDCCP](http://github.com/petar/GoDCCP) that
implements the pure DCCP protocol.

We are looking for contributors!

# About

GoNeedle is written by [Petar Maymounkov](http://pdos.csail.mit.edu/~petar/). 
