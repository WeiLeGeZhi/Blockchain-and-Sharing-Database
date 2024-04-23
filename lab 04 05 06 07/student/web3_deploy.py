from web3 import Web3
from web3.contract import Contract
from solcx import compile_standard, install_solc

# Install Solidity compiler.
_solc_version = "0.7.6"
install_solc(_solc_version)

# ganache 默认rpc url和端口
url = 'http://localhost:7545'
w3 = Web3(Web3.HTTPProvider(url))
if not w3.is_connected():
    raise Exception("Ethereum connection failed!")
eth = w3.eth

# 设置发送交易的账户
eth.default_account = eth.accounts[0]
print("Ethereum({}) connection successful!".format(url))

""" 读取文件中的全部内容 """
def read_file(file) -> str:
    with open(file, 'r', encoding='utf-8') as f:
        return f.read()

""" 部署合约 """
def deploy(contract_name, abi, bytecode, *args) -> Contract:
    # 生成合约，调用其构造函数将其部署至以太坊上
    contract = eth.contract(abi=abi, bytecode=bytecode)
    tx_hash = contract.constructor(*args).transact()
    # 获取交易回执，使用其中的合约地址定位已部署的合约
    tx_receipt = eth.wait_for_transaction_receipt(tx_hash)
    deployed_contract = eth.contract(address=tx_receipt.contractAddress, abi=abi)
    print("Contract '{}' was successfully deployed! Contract address: '{}'".format(contract_name, tx_receipt.contractAddress))
    return deployed_contract

""" 编译合约 """
def compile(contract_name) -> tuple[str, str]:
    contract_sol = contract_name + '.sol'
    contract_content = read_file(contract_sol)
    # Compile smart contract with solcx.
    compiled_result = compile_standard(
        {
            "language": "Solidity",
            "sources": {contract_sol: {"content": contract_content}},
            "settings": {
                "outputSelection": {
                    "*": {"*": ["abi", "evm.bytecode"]}
                }
            },
        },
        solc_version=_solc_version,
    )
    print("Contract '{}' was successfully compiled!".format(contract_name))
    return compiled_result['contracts'][contract_sol][contract_name]['abi'], compiled_result['contracts'][contract_sol][contract_name]['evm']['bytecode']['object']

""" 编译合约并部署 """
def compile_and_deploy(contract_name, *args) -> Contract:
    abi, bytecode = compile(contract_name)
    return deploy(contract_name, abi, bytecode, *args)