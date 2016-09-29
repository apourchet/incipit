#! /usr/bin/ruby

require 'colorize'
require 'json'

require './tools/kubexpose'

module SimpleRpc
    @@port = KubeExpose.expose 'simplerpc'

    def SimpleRpc.test
        print ">> Testing SimpleRpc: "
        
        t0 = SimpleRpc.putKey
        t1 = SimpleRpc.getKey

        if t0 && t1 then print "SUCCESS\n".colorize(:green) else print "FAILED\n".colorize(:red) end
        SimpleRpc.cleanup
    end

    def SimpleRpc.putKey
        uri = URI "http://incipit.machine:#{@@port}/rpc/v1/simplerpc"
        req = Net::HTTP::Post.new uri
        req.body = {'key'=>'foo', 'value'=>'bar'}.to_json
        res = Net::HTTP.start uri.hostname, uri.port do |http|
            http.request(req)
        end
        return (res.code == '200')
    end

    def SimpleRpc.getKey
        uri = URI "http://incipit.machine:#{@@port}/rpc/v1/simplerpc"
        req = Net::HTTP::Get.new uri
        req.body = {'key'=>'foo'}.to_json
        res = Net::HTTP.start uri.hostname, uri.port do |http|
            http.request(req)
        end
        return (res.code == '200')
    end

    def SimpleRpc.cleanup
        KubeExpose.hide 'simplerpc'
    end
end
