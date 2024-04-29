// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.16 <0.9.0;

import "truffle/Assert.sol"; //引入断言
import "truffle/DeployedAddresses.sol"; // 获取被测试合约地址
import "../contracts/Adoption.sol"; //被测试合约

contract TestAdoption {
    Adoption adoption = Adoption(DeployedAddresses.Adoption());

    //领养者测试用例
    function testUserCanAdoptPet() public {
        uint returnId = adoption.adopt(8);

        uint expected = 8;
        Assert.equal(returnId, expected, "Adoption of pet ID 8 should be recorded.");
    }

    //宠物所有者测试用例
    function testGetAdopterAddressByPetId() public {
        //期望领养者地址就是本合约地址，因为交易是测试合约发起的
        address expected = address(this);
        address adopter = adoption.adopters(8);
        Assert.equal(adopter, expected, "Owner of pet Id 8 should be recorded.");
    }

    //测试所有领养者
    function testGetAdopterAddressByPetIdInArray() public {
        //领养者地址就是本合约地址
        address expected = address(this);
        address[16] memory adopters = adoption.getAdopters();
        Assert.equal(adopters[8], expected, "Owner of pet ID 8 should be recorded.");
    }
}