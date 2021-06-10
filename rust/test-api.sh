#!/bin/bash
# executable program `jq` is needed to run the script

echo -e 'Hello'
curl http://127.0.0.1:8090/api/hello -d '{"name": "weaming", "legs": 41}'

# echo -e '\nBIP44'
# curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/44'/1'/0'/0/1\"}" | jq
# echo -e '\nBIP49'
# curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/49'/1'/0'/0/1\"}" | jq
# echo -e '\nBIP84'
# curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/84'/1'/0'/0/1\"}" | jq
# echo -e '\nMultiSig P2SH bad request 1'
# curl -s http://127.0.0.1:8080/api/multisig-p2sh-address -d '{"pubkeys": ["1MdnhGkQ5QNuG3LWfGzuh58AXX8T3XC2sS"], "n": 2}' | jq
# echo -e '\nMultiSig P2SH bad request 2'
# curl -s http://127.0.0.1:8080/api/multisig-p2sh-address -d '{"pubkeys": ["1LqEjz7JefSCJfrhb1ga9KLZmh81iw2QCw", "12W9StPbvenVra5oMmvGxEkJvgTqC6EcYV"], "n": 3}' | jq
# echo -e '\nMultiSig P2SH good request'
# curl -s http://127.0.0.1:8080/api/multisig-p2sh-address -d '{"pubkeys": ["1LqEjz7JefSCJfrhb1ga9KLZmh81iw2QCw", "12W9StPbvenVra5oMmvGxEkJvgTqC6EcYV", "1MdnhGkQ5QNuG3LWfGzuh58AXX8T3XC2sS"], "n": 2}' | jq