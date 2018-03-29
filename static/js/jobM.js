let socket = io();

socket.on('jobUpdate', function (dataStr) {
    const data = JSON.parse(dataStr);
    if (data.id) {
        let tr = $('*[data-id="' + data.id + '"]').closest('tr');
        if(tr && tr.length === 0){
            $('#dataTable').append(createJobTr(data));
            removeEmptyDataTr();
        }else{
            setHtml(tr, 'skuId', data.skuId);
            setHtml(tr, 'quantity', data.quantity);
            setHtml(tr, 'lot', data.lot);
            setHtml(tr, 'outputCount', data.outputCount);
            setHtml(tr, 'stationId', data.stationId);
            setHtml(tr, 'createdAt', getDateFormatFromStr(data.createdAt));

        }
    }
});

function createJobTr(data){
    return '<tr>' +
        `<td data-id="${data.id}">${data.id}</td>` +
        `<td class="skuId">${data.skuId || ''}</td>` +
        `<td class="quantity">${data.quantity || ''}</td>` +
        `<td class="lot">${data.lot || ''}</td>` +
        `<td class="outputCount">${data.outputCount || ''}</td>` +
        `<td class="stationId">${data.stationid || ''}</td>` +
        `<td class="createdAt">${getDateFormatFromStr(data.createdAt, 'YYYY-MM-DD HH:mm:ss')}</td>` +
        '<td>' +
        `<a class="btn btn-primary" href="/jobM/save?id=${data.id}">${$('#strBtnEdit').text()}</a>` +
        `<button class="btn btn-danger" onclick="remove(this, '/jobM', 'id', '${data.id}')">${$('#strDelete').text()}</button>` +
        '</td>' +
        '</tr>';
}


