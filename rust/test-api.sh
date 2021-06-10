#!/bin/bash
# executable program `jq` is needed to run the script

echo -e '\nBIP44 bad 1'
curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/44'/0'/0'/0\"}" | jq
echo -e '\nBIP44 bad 2'
curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/99'/0'/0'/0/1\"}" | jq
echo -e '\nBIP44 bad 3'
curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/-1'/0'/0'/0/1\"}" | jq
echo -e '\nBIP44 bad 4'
curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/44'/2'/0'/0/1\"}" | jq
echo -e '\nBIP44'
curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/44'/0'/0'/0/1\"}" | jq
echo -e '\nBIP49'
curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/49'/0'/0'/0/1\"}" | jq
echo -e '\nBIP84'
curl -s http://127.0.0.1:8080/api/hd-segwit-address -d "{\"path\": \"m/84'/0'/0'/0/1\"}" | jq
# echo -e '\nMultiSig P2SH bad request 1'
# curl -s http://127.0.0.1:8080/api/multisig-p2sh-address -d '{"pubkeys": ["1MdnhGkQ5QNuG3LWfGzuh58AXX8T3XC2sS"], "n": 2}' | jq
# echo -e '\nMultiSig P2SH bad request 2'
# curl -s http://127.0.0.1:8080/api/multisig-p2sh-address -d '{"pubkeys": ["1LqEjz7JefSCJfrhb1ga9KLZmh81iw2QCw", "12W9StPbvenVra5oMmvGxEkJvgTqC6EcYV"], "n": 3}' | jq
# echo -e '\nMultiSig P2SH good request'
# curl -s http://127.0.0.1:8080/api/multisig-p2sh-address -d '{"pubkeys": ["1LqEjz7JefSCJfrhb1ga9KLZmh81iw2QCw", "12W9StPbvenVra5oMmvGxEkJvgTqC6EcYV", "1MdnhGkQ5QNuG3LWfGzuh58AXX8T3XC2sS"], "n": 2}' | jq
