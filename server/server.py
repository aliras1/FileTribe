from flask import Flask, abort, jsonify, request, Response


app = Flask(__name__)
users = {
    #'hali': "K81nggpY96g95wm0blComPnmQw3qmhTMke7llFso9WSQDQZb59Oz9MeO+82gimfr7xO2Q+4Q4SYAGe+wqMScaeOwxEKY/BwUmvv0yJlvuSQnrkHkZJuTTKSVmRt4UrhV"
}

h_to_box_key = {
    #'K81nggpY96g95wm0blComPnmQw3qmhTMke7llFso9WSQDQZb59Oz9MeO+82gimfr7xO2Q+4Q4SYAGe+wqMScaeOwxEKY/BwUmvv0yJlvuSQnrkHkZJuTTKSVmRt4UrhV': 'K81nggpY96g95wm0blComPnmQw3qmhTMke7llFso9WQ='
}
h_to_sign_key = {
    #'K81nggpY96g95wm0blComPnmQw3qmhTMke7llFso9WSQDQZb59Oz9MeO+82gimfr7xO2Q+4Q4SYAGe+wqMScaeOwxEKY/BwUmvv0yJlvuSQnrkHkZJuTTKSVmRt4UrhV': 'kA0GW+fTs/THjvvNoIpn6+8TtkPuEOEmABnvsKjEnGk='
}
h_to_ipfs = {}
messages = {}
groups = {}


@app.route('/is/group/registered/<group_name>', methods=['GET'])
def is_group_registered(group_name):
    if group_name not in groups:
        return Response("false")
    return Response("true")


@app.route('/is/username/registered/<username>', methods=['GET'])
def is_username_registered(username):
    if username not in users:
        return Response("false")
    return Response("true")


@app.route('/get/group/signkey/<group_name>', methods=['GET'])
def get_group_signing_key(group_name):
    if group_name not in groups:
        return Response()
    return groups[group_name]


@app.route('/get/user/publickeyhash/<user>', methods=['GET'])
def get_user_public_key_hash(user):
    if user not in users:
        abort(404)
    return Response(users[user])


@app.route('/get/user/signkey/<user>', methods=['GET'])
def get_user_signing_key(user):
    if user not in users:
        return Response()
    h = users[user]
    return Response(h_to_sign_key[h])


@app.route('/get/user/boxkey/<user>', methods=['GET'])
def get_user_boxing_key(user):
    if user not in users:
        return Response()
    h = users[user]
    return Response(h_to_box_key[h])


@app.route('/get/user/ipfsaddr/<user>', methods=['GET'])
def get_user_ipfs_addr(user):
    if user not in users:
        return Response()
    h = users[user]
    return Response(h_to_ipfs[h])


@app.route('/put/signkey', methods=['POST'])
def put_sign_key():
    data = request.json
    hash = data["hash"]
    vk = data["signkey"]
    print(hash)
    print(vk)
    if hash in h_to_sign_key:
        Response("signing key for hash {} already exists".format(hash))
    h_to_sign_key[hash] = vk
    return Response()


@app.route('/put/boxkey', methods=['POST'])
def put_box_key():
    data = request.json
    hash = data["hash"]
    vk = data["boxkey"]
    print(hash)
    print(vk)
    if hash in h_to_box_key:
        Response("boxing key for hash {} already exists".format(hash))
    h_to_box_key[hash] = vk
    return Response()


@app.route('/put/ipfsaddr', methods=['POST'])
def put_ipfs_addr():
    data = request.json
    hash = data["hash"]
    ipfs = data["ipfsaddr"]
    print(hash)
    print(ipfs)
    if hash in h_to_ipfs:
        Response("ipfs address for hash {} already exists".format(hash))
    h_to_ipfs[hash] = ipfs
    return Response()


@app.route('/register/group/<group_name>', methods=['POST'])
def register_group(group_name):
    data = request.data
    print(data.decode())
    if group_name in groups:
        Response("user already exists")
    groups[group_name] = data.decode()
    return Response()


@app.route('/register/username/<username>', methods=['POST'])
def signup(username):
    data = request.data
    print(data.decode())
    if username in users:
        Response("user already exists")
    users[username] = data.decode()
    return Response()


@app.route('/send/message', methods=['POST'])
def send_message():
    data = request.json
    if data["to"] in messages:
        messages[data["to"]] += [{"from": data["from"], "type": data["type"], "message": data["message"]}]
    else:
        messages[data["to"]] = [{"from": data["from"], "type": data["type"], "message": data["message"]}]
    print(messages)
    return Response()


@app.route('/get/messages/<username>', methods=['GET'])
def get_messages(username):
    if username not in messages:
        return Response(jsonify([]).data)
    else:
        msgs = messages[username]
        del messages[username]
        return Response(jsonify(msgs).data)


if __name__ == "__main__":
    import os
    os.environ['NLS_LANG'] = '.UTF8'
    app.run(debug=True, host="0.0.0.0", port=6000)
