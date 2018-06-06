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
        username = "Alice"
        await dipfshare.registerUser(
            username,
            "0x0000000000000000000000000000000000000000000000000000000000000001",
            "0x0000000000000000000000000000000000000000000000000000000000000002",
            "/ipfs/23h1kjh1329sdfhsdk2323rfsdf"
        );

        const registered = await dipfshare.isUserRegistered(username)
        assert(registered === true);

        const user = await dipfshare.getUser(username)
        assert(user[0] === creator)
        assert(user[1] === "0x0000000000000000000000000000000000000000000000000000000000000001",)
        assert(user[2] === "0x0000000000000000000000000000000000000000000000000000000000000002",)
        assert(user[3] === "/ipfs/23h1kjh1329sdfhsdk2323rfsdf")
    });
});