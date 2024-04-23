// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.15 <0.9.0;
pragma abicoder v2;

contract StudentContract {
    struct Student {
        uint256 id;
        string name;
        string sex;
        uint256 age;
        string dept;
    }

    address admin;
    Student[] students;
    uint256[] ids;
    uint256 count = 0;
    mapping(uint256 => uint256) indexMapping; // id映射index
    mapping(uint256 => bool) isExistMapping;

    constructor() {
        admin = msg.sender;
    }

    function insert(
        uint256 _id,
        string memory _name,
        string memory _sex,
        uint256 _age,
        string memory _dept
    ) public {
        // TODO:插入一条学生记录

        emit Insert(_id);
    }

    event Insert(uint256 id);

    function exist_by_id(uint256 _id) public view returns (bool isExist) {
        // TODO:查找系统中是否存在某个学号
    }

    function select_count() public view returns (uint256 _count) {
        // TODO:查找系统中的学生数量
    }

    function select_all_id() public view returns (uint256[] memory _ids) {
        // TODO:查找系统中所有的学号
    }

    function select_id(uint256 _id) public view returns (Student memory) {
        // TODO:查找指定学号的学生信息
    }

    function delete_by_id(uint256 _id) public {
        // TODO:删除指定学号的学生信息
    }
}
