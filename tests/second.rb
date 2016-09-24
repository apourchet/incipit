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
sleep(0.5)
`curl -s http://dummy.machine:#{port}/ -k`
if $?.success? then print "SUCCESS\n".colorize(:green) else print "FAILED\n".colorize(:red) end
KubeExpose.hide('hermes')