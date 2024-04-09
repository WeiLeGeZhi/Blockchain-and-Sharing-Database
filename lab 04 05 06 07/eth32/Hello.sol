// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.16 <0.9.0;
pragma solidity ^0.8.3;

contract Hello {
    string greeting;

    constructor(string memory _greeting) {
        greeting = _greeting;
    }

    function say() public view returns (string memory) {
        return greeting;
    }
}