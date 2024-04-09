// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.16 <0.9.0;

contract Simple {
    uint data;

    function set(uint _data) public {
        data = _data;
    }

    function get() public view returns (uint) {
        return data;
    }

    function add(uint x) public {
        data += x;
    }

    function fib(uint x) public pure returns (uint) {
        uint a = 0;
        uint b = 1;

        for (uint i=0; i<x; ++i) {
            uint c = b;
            b = a + b;
            a = c;
        }

        return a;
    }

    function hash(string memory str) public pure returns (uint) {
        uint res = uint(keccak256(abi.encodePacked(str)));
        return res;
    }
}