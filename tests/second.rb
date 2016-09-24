#! /usr/bin/ruby
require 'colorize'
require './tools/kubexpose'

print ">> Testing hellogo: "
`curl -s https://dummy.machine:30002/hellogo -k`
if $?.success? then print "SUCCESS\n".colorize(:green) else print "FAILED\n".colorize(:red) end

print ">> Testing hellonode: "
`curl -s https://dummy.machine:30002/hellonode -k`
if $?.success? then print "SUCCESS\n".colorize(:green) else print "FAILED\n".colorize(:red) end

print ">> Testing hermes: "
port = KubeExpose.expose('hermes')
puts `curl -s http://dummy.machine:#{port}/rpc/v1/hermes -k`
if $?.success? then print "SUCCESS\n".colorize(:green) else print "FAILED\n".colorize(:red) end
KubeExpose.hide('hermes')
