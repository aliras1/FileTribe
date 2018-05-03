from flask import Flask, abort, jsonify, request, Response
import nacl.encoding
import nacl.signing
import nacl.exceptions
import base64


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


@app.route('/get/group/members/<group_name>', methods=['GET'])
def get_group_signing_key(group_name):
    if group_name not in groups:
        return Response()
    print(groups)
    return Response(jsonify(groups[group_name]["members"]).data)


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


@app.route('/register/group', methods=['POST'])
def register_group():
    data = request.json
    group_name = data["groupname"]
    owner = data["owner"]
    state = data["state"]
    if group_name in groups:
        Response("group already exists")
    groups[group_name] = {"members": [owner], "state": [state], "last_op": [None]}
    print(groups[group_name])
    return Response()


@app.route('/get/group/state/<group_name>', methods=['GET'])
def get_group_state(group_name):
    if group_name not in groups:
        print("group {} does not exist: get_group_state()".format(group_name))
        return Response("group does not exist")
    return Response(groups[group_name]["state"][-1])


@app.route('/get/group/prev/state/<group_name>/<state>', methods=['GET'])
def get_group_prev_state(group_name, state):
    if group_name not in groups:
        print("group {} does not exist: get_group_state()".format(group_name))
        return Response("group does not exist")
    if state not in groups[group_name]["state"]:
        print("state not in group state history: get_group_prev_state()")
    for i in range(1, len(groups[group_name]["state"])):
        if groups[group_name]["state"][i] == state:
            return Response(groups[group_name]["state"][i-1])
    return Response()


def verify(signed, verify_key):
    try:
        verify_key.verify(signed)
    except nacl.exceptions.BadSignatureError:
        return False
    return True


def verify_transaction(transaction):
    for signed_by in transaction["signed_by"]:
        username = signed_by["username"]
        signature_base64 = signed_by["signature"]
        signature = bytearray(base64.b64decode(signature_base64))
        for b in base64.b64decode(transaction["hash"]):
            signature.append(b)
        signed = [b for b in signature]
        h = users[username]
        verify_key = nacl.signing.VerifyKey(h_to_sign_key[h], encoder=nacl.encoding.Base64Encoder) 

        if not verify(signed, verify_key):
            return False
    return True



@app.route('/group/invite/<group_name>', methods=['POST'])
def group_invite(group_name):
    if group_name not in groups:
        print("group {} does not exists: group_invite()".format(group_name))
        return Response()
    transaction = request.json
    if not verify_transaction(transaction):
        return Response()
    groups[group_name]["state"] += [transaction["hash"]]
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
