require 'net/http'
require 'json'

require './tests/hellogo/unit'
require './tests/hellonode/unit'
require './tests/simplerpc/unit'

HelloGo.test()
HelloNode.test()
SimpleRpc.test()
