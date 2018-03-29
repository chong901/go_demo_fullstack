let socket = io.connect();

socket.on('jobUpdate', function(data){
  if(data.masterId){
    var tr = $('*[data-masterId="' + data.masterId + '"]').closest('tr');
    setHtml(tr, 'stationId', data.stationId);
    setHtml(tr, 'workerId', data.workerId);
    setHtml(tr, 'outputCount', data.outputCount);
    setHtml(tr, 'quantity', data.quantity);
  }
});
