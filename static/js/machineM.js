let socket = io.connect();

socket.on('machineUpdate', function (dataStr) {
    const data = JSON.parse(dataStr);
    if (data.id) {
        let tr = $('*[data-id="' + data.id + '"]').closest('tr');
        setHtml(tr, 'name', data.name);
        setHtml(tr, 'type', data.type);
        setHtml(tr, 'department', data.department);
    }
});
