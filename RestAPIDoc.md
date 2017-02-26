#RESTAPI Documentation  

##Endpoints  

###Sign-up  

* /signup
   
   A GET request will return ```{"id":"example"}```. That will be your client_id from now on, and will be saved in the server database.


###Channels  
####All of the following endpoints require ?client_id as an url parameter.

* /channels/list

   Accepts a GET request.
   Returns a list of json encoded channels in this form ```{ [{"id":"abcd", name:"name"}, ...]}```  
  
  
* /channels/create  

   Creates a new channel.   
   Accepts a POST request of the form ```{"name":"new channel name here"}```  
  
  
* /channels/delete

   Deletes a channel.  
   Accepts a POST request of the form ```{"id":"channel id here"}```




###User

* /user/name

   Updates the user nickname. It can also return the current user nickname.
   Accepts a POST request with ```{"name":"example"}```, to update the name.
   Accepts a GET request, returns ```{"name":"example"}```.
   
* /user/channels

   Returns a JSON array of channels the user is subscribed to. 
   Accepts a GET request returns ```{ [{"id":"abcd", name:"name"}, ...]}```
   


