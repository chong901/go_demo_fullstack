function editFunction(action, dom, id) {
    let nameMap;
    $.ajax({
        url: '/api/v1/function/funcDict'
    }).then((r) => {
        nameMap = r.data.reduce((pre, cur) => {
            if (!cur) return pre;
            pre[cur.id] = cur.name;
            return pre;
        }, {});
        if (id) {
            return $.ajax(`/api/v1/function?id=${id}`);
        } else {
            return {data: {}};
        }
    }).then((r) => {
        createDialog(action, dom, r.data, nameMap);
    }).fail(alertErr);
}

function createDialog(action, dom, funcData, namesMap) {
    BootstrapDialog.show({
        title: $('#strBtnEdit').text(),
        closable: false,
        message: createForm(funcData, namesMap),
        onshown: function (bd) {
            $(".selectpicker").selectpicker('refresh');
            $("#functionForm").on('submit', function (ev) {
                ev.preventDefault();

                const data = getObjFromSerialize($(this).serializeArray(), ['id', 'parent', 'orderNum'], null, ['isMenu']);
                $.ajax({
                    url: '/function',
                    method: 'POST',
                    dataType: 'json',
                    data: JSON.stringify(data)
                }).done((r) => {
                    const data = r.data;
                    switch (action) {
                        case 'edit':
                            const $tr = $(dom).closest('tr');
                            setHtmlAcceptEmpty($tr, 'parentName', namesMap[data.parent]);
                            setHtml($tr, 'name', data.name);
                            setHtml($tr, 'uri', data.uri);
                            setHtml($tr, 'method', data.method);
                            setHtml($tr, 'orderNum', data.orderNum);
                            setHtml($tr, 'isMenu', data.isMenu ? 'T' : 'F');
                            break;
                        case 'add':
                            $('#dataTable').append(createTr(data, namesMap));
                            removeEmptyDataTr();
                            break;
                    }
                    bd.close();
                }).fail(alertErr);
            });
        },
        buttons: [
            {
                label: $('#strBtnConfirm').text(),

                action: function () {
                    $("#functionFormSubmit").click();
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (bd) {
                    bd.close();
                }
            }
        ]
    });
}

function createForm(funcData, namesMap) {
    const {id, parent, name, uri, method, orderNum, isMenu: isMenu} = funcData;
    let html = '<form id="functionForm">' +
        `<input type="hidden" name="id" value="${id || ''}" />` +
        '<div class="form-group">' +
        `<label for="parent">${$('#strParent').text()}</label><br/>` +
        createNameSelection(namesMap, parent) +
        '</div>' +
        '<div class="form-group">' +
        `<label for="name">${$('.strName').text()}</label>` +
        '<input id="name" class="form-control" name="name" required value="' + (name || '') + '"/>' +
        '</div>' +
        '<div class="form-group">' +
        `<label for="uri">${$('.strURI').text()}</label>` +
        '<input id="uri" class="form-control" name="uri" value="' + (uri || '') + '"/>' +
        '</div>' +
        `<label for="method">${$('.strMethod').text()}</label>` +
        createMethodSelect(method) +
        '</div>' +
        '<div class="form-group">' +
        `<label for="order">${$('.strOrder').text()}</label>` +
        `<input id="order" class="form-control" name="orderNum" type="number" value="${orderNum || 0}"/>` +
        '</div>' +
        '<div class="form-group">' +
        '<label class="form-group">' +
        `<label>${$('.strIsMenu').text()}</label><br/>` +
        `<label class="radio-inline"><input name="isMenu" type="radio" value="1" ${isMenu ? 'checked' : ''}/>${$('#strYes').text()}</label>` +
        `<label class="radio-inline"><input name="isMenu" type="radio" value="0" ${isMenu ? '' : 'checked'}/>${$('#strNo').text()}</label>` +
        '</div>' +
        '<input id="functionFormSubmit" type="submit" style="display: none"/>' +
        '</form>';
    return html;
}

function createNameSelection(nameMap, val) {
    let html = '<select id="parent" name="parent" class="selectpicker" data-live-search="true"><option></option>';
    for (let k in nameMap) {
        html += `<option value="${k || ''}" ${parseInt(k) === parseInt(val) ? 'selected' : ''}>${nameMap[k]}</option>`;
    }
    return html + '</select>';
}

function createTr(data, nameMap) {
    let html = '<tr>' +
        '<td class="parentName">' + (nameMap[data.parent] || '') + '</td>' +
        '<td class="name">' + data.name + '</td>' +
        '<td class="uri">' + data.uri + '</td>' +
        '<td class="method">' + data.method + '</td>' +
        '<td class="isMenu">' + (data.isMenu ? 'T' : 'F') + '</td>' +
        '<td class="orderNum">' + data.orderNum + '</td>' +
        '<td>' +
        `<input class="parent" type="hidden" value="${data.parent}"/>` +
        `<input class="btn btn-primary updateBtn" type="button" value="${$('#strBtnEdit').text()}" onclick="editFunction('edit', this, '${data.id}')"/> ` +
        `<input class="btn btn-danger" type="button" value="${$('#strDelete').text()}" onclick="remove(this, '/function', 'id', '${data.id}')"/>` +
        '</td>' +
        '</tr>';
    return html;
}

function createMethodSelect(val) {
    let html = '<select id="method" name="method" class="form-control">' +
        '<option></option>' +
        '<option value="GET" ' + (val === 'GET' ? 'selected' : '') + '>GET</option>' +
        '<option value="POST" ' + (val === 'POST' ? 'selected' : '') + '>POST</option>' +
        '<option value="PUT" ' + (val === 'PUT' ? 'selected' : '') + '>PUT</option>' +
        '<option value="DELETE" ' + (val === 'DELETE' ? 'selected' : '') + '>DELETE</option>' +
        '</select>';
    return html;
}