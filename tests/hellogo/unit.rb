#! /usr/bin/ruby

require 'colorize'
require 'json'
require 'net/http'
require 'openssl'

require './tools/kubexpose'

module HelloGo
    @@port = 30002

    def HelloGo.test
        print ">> Testing HelloGo: "
        
        t0 = HelloGo.hello

        if t0 then print "SUCCESS\n".colorize(:green) else print "FAILED\n".colorize(:red) end
    end

    def HelloGo.hello
        uri = URI "https://dummy.machine:#{@@port}/hellogo"
        req = Net::HTTP::Post.new uri
        options = {:use_ssl => true, :verify_mode => OpenSSL::SSL::VERIFY_NONE}
        res = Net::HTTP.start(uri.hostname, uri.port, options) do |http|
            http.request(req)
        end
        return (res.code == '200')
    end
end
