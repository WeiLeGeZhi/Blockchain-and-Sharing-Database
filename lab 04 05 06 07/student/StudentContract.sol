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
        require(msg.sender == admin, "Only admin can insert students.");

        if (!isExistMapping[_id]) {
            Student memory newStudent = Student(_id, _name, _sex, _age, _dept);
            students.push(newStudent);
            ids.push(_id);
            indexMapping[_id] = students.length - 1;
            isExistMapping[_id] = true;
            count += 1;
            emit Insert(_id);
        }
    }

    event Insert(uint256 id);

    function exist_by_id(uint256 _id) public view returns (bool isExist) {
        // TODO:查找系统中是否存在某个学号
        return isExistMapping[_id];
    }

    function select_count() public view returns (uint256 _count) {
        // TODO:查找系统中的学生数量
        return students.length;
    }

    function select_all_id() public view returns (uint256[] memory _ids) {
        // TODO:查找系统中所有的学号
        return ids;
    }

    function select_id(uint256 _id) public view returns (Student memory) {
        // TODO:查找指定学号的学生信息
        require(isExistMapping[_id], "Student does not exist.");
        return students[indexMapping[_id]];
    }

    function delete_by_id(uint256 _id) public {
        // TODO:删除指定学号的学生信息
        require(msg.sender == admin, "Only admin can delete students.");

        if (isExistMapping[_id]) {
            uint256 index = indexMapping[_id];
            uint256 lastIndex = students.length - 1;

            // Swap the element to delete with the last element
            students[index] = students[lastIndex];
            ids[index] = ids[lastIndex];
            indexMapping[ids[lastIndex]] = index;

            // Remove the last element
            students.pop();
            ids.pop();
    
            // Clear mappings
            delete indexMapping[_id];
            isExistMapping[_id] = false;
            count -= 1;
        }
    }

    function get_id_by_min_age() public view returns (uint256 id) {
        require(count > 0, "No students in the system.");

        uint256 minAge = students[0].age;
        id = students[0].id;

        for (uint256 i = 1; i < count; i++) {
            if (students[i].age < minAge) {
                minAge = students[i].age;
                id = students[i].id;
            }
        }

        return id;
    }

    function update_dept_by_id(uint256 _id, string memory _dept) public {
        require(msg.sender == admin, "Only admin can update student's dept.");
        require(isExistMapping[_id], "Student does not exist.");

        uint256 index = indexMapping[_id];
        students[index].dept = _dept;

        emit Update(_id, _dept);
    }

    event Update(uint256 id, string dept);
    
}
