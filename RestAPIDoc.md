#RESTAPI Documentation  

##Endpoints  


###Channels  

* /channels/list

   Accepts a POST request of the form ```{ "client_id":"id_here"}```
   Returns a list of json encoded channels in this form ```{ [{"id":"abcd", name:"name"}, ...]}```  
  
  
* /channels/create  

   Creates a new channel.   
   Accepts a POST request of the form ```{"client_id":"client id here", "name":"new channel name here"}```  
  
  
* /channels/delete

   Deletes a channel.  
   Accepts a POST request of the form ```{"client_id":"client id here", "id":"channel id here"}```

###Sign-up  

* /signup
   
   A GET request will return ```{"id":"example"}```. That will be your client_id from now on, and will be saved in the server database.


###User

* /user/name
   Updates the user nickname. It can also return the current user nickname if the name parameter in the POST request is empty.
   Accepts a POST request with ```{ "client_id":"id_here", "name":"example"}```, returns ```{"name":"name_here"}```
   
* /user/channels
   Returns a JSON array of channels the user is subscribed to. 
   Accepts a POST request with ```{ "client_id":"id_here"}```, returns ```{ [{"id":"abcd", name:"name"}, ...]}```
   


