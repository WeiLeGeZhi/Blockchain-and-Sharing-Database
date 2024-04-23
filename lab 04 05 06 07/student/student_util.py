# 生成随机的学生数据进行测试
import random
import string

depts = [
    '教育学院',
    '计算机学院',
    '数据学院',
    '软件学院',
    '数学学院',
    '物理学院',
    '化学学院',
    '艺术学院',
    '经济学院',
    '哲学学院',
    '生物学院',
    '音乐学院',
    '外语学院',
    '历史学院',
    '管理学院'
]

def random_name():
    size = random.randint(5, 8)
    first_char = random.choice(string.ascii_uppercase)
    remain = ''.join(random.choices(string.ascii_lowercase, k=size-1))
    return first_char + remain

def random_sex():
    return random.choice(['male', 'female'])

def random_age():
    return random.randint(18, 24)

def random_dept():
    return random.choice(depts)

def random_student_list(num):
    students = dict()
    ids = random.sample(range(0, int(10e9)), num)
    for id in ids:
        # 学生元组形式: (id, name, sex, age, dept)
        student = (id, random_name(), random_sex(), random_age(), random_dept())
        students[id] = student
        # print(student)
    return students


def print_error(msg):
    print("\033[31m{}\033[0m".format(msg))

def print_success(msg):
    print("\033[33m{}\033[0m".format(msg))

def print_info(msg):
    print("\033[32m{}\033[0m".format(msg))


# print_success("success!")
# print_error("error!")
# print_info("info")