for (var i = 0; i < 1; i++) {
  console.log('\t\t' + (i + 1) + ' - Transactions')
    web3.eth.sendTransaction({
        from: web3.eth.accounts[3],
        to: web3.eth.accounts[2],
        value: 5,
        gas: 50000,
        gasPrice: 20
    });
    web3.eth.sendTransaction({
        from: web3.eth.accounts[3],
        to: web3.eth.accounts[2],
        value: 291,
        gas: 50000,
        gasPrice: 20
    });

    web3.eth.sendTransaction({
        from: web3.eth.accounts[3],
        to: web3.eth.accounts[1],
        value: 53039,
        gas: 50000,
        gasPrice: 20
    });
}