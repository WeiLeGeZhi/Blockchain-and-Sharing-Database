    // ethereum.request({ method: 'eth_requestAccounts' }).then((accounts) => {
    //     let account = accounts[0]
    //     ArtComissionContract.deployed().then((instance) => {
    //         console.log(instance)
    //         return instance.post('001', 'author', 'you', '10000', 'noninfo', { from: account })
    //         return instance.post(
    //             post_comission.id, post_comission.author, post_comission.title, 
    //             post_comission.price, post_comission.info, 
    //             { from: account }
    //         )
    //     }).then((response) => {
    //         alert('success', response)
    //         loadPostedData()
    //     }).catch((err) => {
    //         console.log(err.message);
    //     })
    // })