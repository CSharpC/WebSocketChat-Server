#RESTAPI Documentation  

##Endpoints  


###Channels  

* /channels/list

   Returns a list of json encoded channels in this form ```{ [{"id":"abcd", name:"name"}, ...]}```  

  
* /channels/create

   Creates a new channel  
   Accepts a POST request of the form ```{"client_id":"client id here", "name":"new channel name here"}```  
  
  
* /channels/delete

   Deletes a channel.  
   Accepts a POST request of the form `````{"client_id":"client id here", "id":"channel id here"}```  



