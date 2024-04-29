from web3 import Web3
import os
import uuid
import argparse

url                     = "http://localhost:7545"               # 以太坊测试链 rpc 连接端口
contract_address_file   = 'contract_address.txt'                # 合约地址保存文件
abi_file                = "AeroplaneChess/AeroplaneChess_sol_AeroplaneChess.abi"   # abi 文件
bytecode_file           = "AeroplaneChess/AeroplaneChess_sol_AeroplaneChess.bin"   # 字节码文件

# 连接测试链
w3  = Web3(Web3.HTTPProvider(url))   
eth = w3.eth
assert w3.is_connected(), "Connect to ethereum failed!"
print("eth connect:", w3.is_connected())

destination     = 20   # 目的地位置
dice            = 6    # 骰子
player_num      = 3    # 玩家数量
player_name     = ''
account_id      = -1

""" 进行游戏 """
def play():
    # 获取游戏合约
    global abi_file, bytecode_file
    abi                 = get_abi_from_file(abi_file)
    bytecode            = get_bytecode_from_file(bytecode_file)
    deployed_contract   = get_deployed_contract(abi, bytecode)

    # 打印游戏信息
    game_info = deployed_contract.functions.getGameInfo().call()
    print("game info: destination={}, dice={}, gameover={}, player_num={}".
    format(game_info[0], game_info[1], bool(game_info[2]), game_info[4]))

    # 如果之前已经参加，则继续游戏，合约保存了玩家参加的状态
    if deployed_contract.functions.isJoin().call():
        player = deployed_contract.functions.getPlayer().call()
        print(player_name, "continue this game, player status:", player)
    # 否则参加游戏
    else:
        deployed_contract.functions.join(player_name).transact()
        player = deployed_contract.functions.getPlayer().call()
        print(player_name, "joined this game, player status:", player)

    # 游戏未结束就不断循环
    while not deployed_contract.functions.isOver().call():
        if deployed_contract.functions.isStart().call() and deployed_contract.functions.isMyRound().call():
            # 使用uuid每次生成不同的字符串作为参数，以获取随机骰子
            uuid_str = str(uuid.uuid4())
            print()
            print(player_name, "playing, input:", uuid_str)
            last_status = deployed_contract.functions.getPlayer().call()
            # 由于同步问题，可能执行到此时，已有其他玩家到达终点，游戏结束，这里会抛出异常, 故直接退出循环终止
            try:
                deployed_contract.functions.play(uuid_str).transact()
            except Exception as e:
                print(e)
                show_global_status(deployed_contract)   # 输出结束状态
                break
            # 打印本次执行后的状态
            deployed_contract.functions.getPlayer().transact()
            cur_status = deployed_contract.functions.getPlayer().call()
            print("round:", cur_status[5], "\trandom dice:", cur_status[6], "\tstep:", last_status[2], "->", cur_status[2])
            show_global_status(deployed_contract)
        else:
            # print("The game is not started yet, please wait...")
            pass

    # 打印游戏胜利者
    winner = deployed_contract.functions.getWinner().call()
    print("\ngameover! The winner is", winner)

    # 删除合约地址文件夹
    try: os.remove(contract_address_file) 
    # except Exception as e: print(e)
    except: pass

""" 打印棋盘和所有参与者的状态 """
def show_global_status(deployed_contract):
    players = deployed_contract.functions.getPlayers().call() # 获取所有玩家状态
    board = [[] for _ in range(destination + 1)]
    for player in players:
        board[player[2]].append(player)
    # 打印棋盘
    print('-' * 20 + 'game board' + '-' * 20)
    for idx in range(len(board)):
        print(idx, ':', end=' ')
        for player in board[idx]:
            name, direction, step, round = player[1], player[3], player[4], player[5]
            # 使用符号表示当前玩家向上还是向下
            print('[', name, ('V' if direction == 1 else '^') + str(step), 'r' + str(round), ']', end = ' ')
        print()
    print('-' * 20 + 'game board' + '-' * 20)

""" 设置调用合约、发送交易的账户 """
def set_default_account():
    eth.default_account = eth.accounts[account_id]
    eth.account = eth.default_account

""" 从文件中获取abi """
def get_abi_from_file(file): 
    with open(file, 'r') as f: print("Get abi success!"); return f.read() 

""" 从文件中获取字节码 """
def get_bytecode_from_file(file):
    with open(file, 'r') as f: print("Get bytecode success!"); return "0x" + f.read()

""" 部署合约 """
def deploy_contract(abi, bytecode):
    contract = eth.contract(abi=abi, bytecode=bytecode) # 创建合约
    tx_hash = contract.constructor(destination, dice, player_num).transact() # 部署合约（发送构造函数的交易，需相对应合约中的参数）
    tx_receipt = eth.wait_for_transaction_receipt(tx_hash) # 等待交易回执
    print("contract address:", tx_receipt.contractAddress) # 合约地址
    # 保存合约
    global contract_address_file
    with open(contract_address_file, "w") as f:
        f.write(tx_receipt.contractAddress)
    # 通过地址获取已部署合约
    deployed_contract = eth.contract(address=tx_receipt.contractAddress, abi=abi)
    print("Deploy contract success!")
    return deployed_contract

""" 获取部署合约，如果本地已保存合约地址，则调用该地址的合约，否则重新创建一个新的合约 """
def get_deployed_contract(abi, bytecode):
    # 如果本地已保存合约地址，则调用该地址的合约
    if os.path.exists(contract_address_file):
        print("Get deployed contract...")
        with open(contract_address_file, "r") as f:
            contract_address = f.read()
        print("Get deployed contract success!")
        print("contract address:", contract_address)
        deployed_contract = eth.contract(address=contract_address, abi=abi)
        return deployed_contract
    # 否则重新创建一个新的合约
    else:
        return deploy_contract(abi, bytecode)
    
if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Aeroplane Chess Game")
    parser.add_argument("-n", "--name", required=True, type=str, help="player name")
    parser.add_argument("-i", "--id", required=True, type=int, help="account id")
    args = parser.parse_args()
    if args.name:
        player_name = args.name
    if args.id:
        account_id = int(args.id)

    assert player_name != "" and account_id > 0, 'Need player name and account id, account id must be greater than 0!'
    set_default_account()
    play()