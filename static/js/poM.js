let socket = io();

socket.on('poUpdate', function (dataStr) {
    const data = JSON.parse(dataStr);
    if (data.id) {
        const tr = $('*[data-id="' + data.id+ '"]').closest('tr');
        if (tr && tr.length === 0) {
            $('#dataTable').append(createTr(data));
        } else {
            setHtml(tr, 'status', data.status);
            setHtml(tr, 'clientId', data.clientId);
            setHtml(tr, 'skuList', createSkuData(data.skuList));
            setHtml(tr, 'due', getDateFormatFromStr(data.due, 'YYYY-MM-DD'));
            setHtml(tr, 'status', data.status);
            setHtml(tr, 'createAt', getDateFormatFromStr(data.createAt, 'YYYY-MM-DD HH:mm:ss'));
        }
    }
});

function createTr(data) {
    let html = '<tr>';
    html += `<td data-id="${data.id}">${data.id}</td>`;
    html += `<td class="clientId">${data.clientId || ''}</td>`;
    html += `<td class="skuList">${createSkuData(data.skuList)}</td>`;
    html += `<td class="due">${getDateFormatFromStr(data.due, 'YYYY-MM-DD')}</td>`;
    html += `<td class="status">${data.status || ''}</td>`;
    html += `<td class="createdAt">${getDateFormatFromStr(data.createdAt, 'YYYY-MM-DD HH:mm:ss')}</td>`;
    html += `<td>`;
    html += `<a class="btn btn-primary" href="/poM/save?id=${data.id}">${$('#strBtnEdit').text()}</a> `;
    html += `<button class="btn btn-danger" onclick="remove(this, '/poM', 'id', '${data.id}')">${$('#strDelete').text()}</button> `;
    html += `</td>`;
    html += '</tr>';
    return html;
}

function createSkuData(skuData) {
    if (!skuData) {
        return '';
    }

    return skuData.reduce(function (pre, cur) {
        if (!cur.skuId) {
            return pre;
        }
        return `${pre}${cur.skuId} - ${cur.quantity || '0'} ${cur.unit || ''}<br>`;
    }, '');
}
