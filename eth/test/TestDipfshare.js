const Dipfshare = artifacts.require('Dipfshare');

contract('Dipfshare', accounts => {
    let dipfshare;
    const creator = accounts[0];

    beforeEach(async function () {
        dipfshare = await Dipfshare.new({ from: creator });
    });

    it('check owner', async function () {
        const owner = await dipfshare.owner();
        assert(owner === creator);
    });

    it('register user', async function () {
        id = "0x0000000000000000000000000000000000000000000000000000000000000001"
        username = "Alice"
        publicKey = "0x0000000000000000000000000000000000000000000000000000000000000002"
        verifyKey = "0x0000000000000000000000000000000000000000000000000000000000000003"
        ipfsAddress = "rweiro34n3453uz4grsrd"

        await dipfshare.registerUser(
            id,
            username,
            publicKey,
            verifyKey,
            ipfsAddress
        );

        const registered = await dipfshare.isUserRegistered(id)
        assert(registered === true);

        const user = await dipfshare.getUser(id)
        assert(user[0] === creator)
        assert(user[1] === username)
        assert(user[2] === publicKey)
        assert(user[3] === verifyKey)
        assert(user[4] === ipfsAddress)
    });
});