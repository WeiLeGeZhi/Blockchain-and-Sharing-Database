from student_util import *
from web3_deploy import *

point = 0 # 测试分值

def add_point(val, msg):
    global point
    point += val
    print_success("{}\n\t\t==> Get {} point, current point: {}".format(msg, val, point))

# 合约部署者即为管理员，后面有对管理员的测试，只有管理员才可以进行插入操作
print_info("\n[Test compile_and_deploy]")
try:
    student_contract = compile_and_deploy("StudentContract") # 编译并部署合约
    add_point(20, "compile_and_deploy StudentContract success!")
except:
    print_error("compile_and_deploy StudentContract error!")

num = 10
students = random_student_list(num) # 获取随机生成学生列表
actual_num = 0


# 测试插入
print_info("\n[Test insert and logs]")
for id, student in students.items():
    try:
        # 构造函数通过rpc向以太坊发送该函数交易
        # *为解包操作，即传入的参数为: id, name, sex, age, dept
        tx_hash = student_contract.functions.insert(*student).transact()
        tx_receipt = eth.wait_for_transaction_receipt(tx_hash)
        if student_contract.functions.exist_by_id(id).call():
            add_point(5, "insert {} success!".format(student))
            actual_num += 1
        else:
            print_error("insert {} failed!".format(student))
        # 测试获取日志
        try:
            insert_event = student_contract.events.Insert.get_logs()[0]
            id = insert_event.get('args').get('id')
            if student[0] == id:
                add_point(5, "emit event success! event args: {}".format(insert_event.get('args')))
            else:
                raise Exception
        except Exception as e:
            print_error("error: {}\nget event error! id = {}".format(e, student[0]))
        # 测试插入操作的幂等性
        tx_hash = student_contract.functions.insert(*student).transact()
        tx_receipt = eth.wait_for_transaction_receipt(tx_hash)
    except Exception as e:
        print_error("error: {}\ninsert error! student: {}".format(e, student))

if actual_num == num:
    add_point(10, "actual_num(={}) == num(={}) success!".format(actual_num, num))
else:
    print_error("actual_num(={}) == num(={}) failed!".format(actual_num, num))


# 测试查询
# 测试 select count，返回合约中学生的总数
print_info("\n[Test select_count]")
try:
    select_count = student_contract.functions.select_count().call()
    if actual_num == select_count:
        add_point(10, "actual_num(={}) == select_count(={}) success!".format(actual_num, select_count))
    else:
        print_error("actual_num(={}) == select_count(={}) failed!".format(actual_num, select_count))
except Exception as e:
    print_error("error: {}\nselect_count error!".format(e))

# 测试 select all id，返回合约中所有id构成的数组
print_info("\n[Test select_all_id]")
try:
    all_id = student_contract.functions.select_all_id().call()
    if set(all_id) == set(students.keys()):
        add_point(20, "select_all_id success, all_id: {}".format(all_id))
    else:
        print_error("select_all_id failed: {} != {}".format(set(all_id), set(students.keys())))
except Exception as e:
    print_error("error: {}\nselect_all_id error!".format(e))

# 测试 select id，返回指定id的学生的信息元组
print_info("\n[Test select_id]")
for id in all_id:
    try:
        student = student_contract.functions.select_id(id).call()
        student = tuple(student)
        add_point(2, "select_id(id={}) success!".format(id))
        if student == students[id]:
            add_point(5, "{} == {} success!".format(student, students[id]))
        else:
            print_error("{} == {} failed!".format(student, students[id]))
    except Exception as e:
         print_error("error: {}\nselect_id(id={}) failed!".format(e, id))

# 测试删除
print_info("\n[Test delete_by_id and exist_by_id]")
delete_id_list = random.sample(all_id, len(all_id) // 2)
print("delete_id_list:", delete_id_list)
for id in delete_id_list:
    try:
        tx_hash = student_contract.functions.delete_by_id(id).transact()
        tx_receipt = eth.wait_for_transaction_receipt(tx_hash)
        # 测试删除的幂等性
        tx_hash = student_contract.functions.delete_by_id(id).transact()
        tx_receipt = eth.wait_for_transaction_receipt(tx_hash)
        if not student_contract.functions.exist_by_id(id).call():
            add_point(5, "delete(id={}) success!".format(id))
        else:
            print_error("delete(id={}) failed!".format(id))
    except:
        print_error("delete(id={}) failed!".format(id))

# 再次测试 select count
print_info("\n[Test select_count after delete]")
actual_num -= len(delete_id_list)
try:
    select_count = student_contract.functions.select_count().call()
    if actual_num == select_count:
        add_point(10, "actual_num(={}) == select_count(={}) success!".format(actual_num, select_count))
    else:
        print_error("actual_num(={}) == select_count(={}) failed!".format(actual_num, select_count))
except Exception as e:
    print_error("error: {}\nselect_count error!".format(e))

# 再次测试 select all_id
# 测试 select all id，返回合约中所有id构成的数组
print_info("\n[Test select_all_id again]")
try:
    all_id2 = student_contract.functions.select_all_id().call()
    if set(all_id2) == set(all_id) - set(delete_id_list):
        add_point(20, "select_all_id(={}) success!".format(all_id2))
    else:
        print_error("select_all_id failed: {} != {}".format(set(all_id2), set(all_id) - set(delete_id_list)))
except Exception as e:
    print_error("error: {}\nselect_all_id error!".format(e))


# 测试是否只有管理员才可进行插入和删除操作
print_info("\n[Test admin privilege on insert and delete]")
eth.default_account = eth.accounts[1]
try:
    func_tx = student_contract.functions.insert(*student)
    try:
        tx_hash = func_tx.transact()
        tx_receipt = eth.wait_for_transaction_receipt(tx_hash)
        print_error("only admin insert failed!")
    except:
        add_point(10, "only admin insert success!")
except:
    print_error("insert func error!")

try:
    func_tx = student_contract.functions.delete_by_id(0)
    try:
        tx_hash = func_tx.transact()
        tx_receipt = eth.wait_for_transaction_receipt(tx_hash)
        print_error("only admin delete failed!")
    except:
        add_point(10, "only admin delete success!")
except:
    print_error("delete func error!")


# 打印最终点数
print("point: {}".format(point))