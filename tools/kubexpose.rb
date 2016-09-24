#! /usr/bin/ruby

module KubeExpose
    def KubeExpose.newport()
        return 31000 + Random.rand(1000)
    end

    def KubeExpose.expose(servicename, port=newport())
        `kubectl patch service #{servicename} --type='json' \
        -p='[{"op": "replace", "path": "/spec/type", "value": "NodePort"},
            {"op": "replace", "path": "/spec/ports/0/nodePort", "value": "#{port}"}]'`
        return (if $?.success? then port else 0 end)
    end

    def KubeExpose.hide(servicename)
        `kubectl patch service #{servicename} --type='json' \
        -p='[{"op": "replace", "path": "/spec/type", "value": "ClusterIP"},
            {"op": "remove", "path": "/spec/ports/0/nodePort"}]'`
        return $?.success?
    end
end

if __FILE__ == $0
    if ARGV.length != 2 || (ARGV[1] != 'expose' && ARGV[1] != 'hide')
        puts "Usage: kubexpose.rb <servicename> <expose/hide>"
        exit 1
    end
    servicename = ARGV[0]
    action = ARGV[1]

    if action == "expose"
        res = KubeExpose.expose(servicename)
        if res != 0 then puts res end
    else 
        res = KubeExpose.hide(servicename)
    end
end
