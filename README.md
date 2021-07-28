## To host on Google Cloud Functions
(zipped file with structure)
- function.go
- go.mod
- go.sum
- services/  

## Setting Webhook
TELEGRAM_TOKEN=""  
CLOUD_FUNCTION_URL=""  

curl --data "url=$CLOUD_FUNCTION_URL" https://api.telegram.org/bot$TELEGRAM_TOKEN/SetWebhook  `

# Workflow
(Bolded words are user states)
- **Any state**
    - /start, /start@toGoListBot
        - Sends basic info
        - goto **Idle** 
    - /reset, /reset@toGoListBot
        - Sends basic info
        - goto **Idle**
    - /help, /help@toGoListBot
        - Sends info on commands
        - goto **Idle**
    - /query, /query@toGoListBot
        - Prompt for query type
        - goto **QuerySelectType**
    - /additem, /additem@toGoListBot, 
        - If in group chat
            - Redirect to bot's chat, with "/start addItem" as default first message
        - If already in bot's chat
            - Prompt for name of item to add
            - goto **ReadyForNextAction**
    - /start addItem
        - Prompt for name of item to add
        - goto **ReadyForNextAction**
    - /deleteitem, /deleteitem@toGoListBot
        - Prompt for item to delete
        - goto **DeleteSelect**
    - /edititem, /edititem@toGoListBot
        - If in group chat
            - Redirect to bot's chat, with "/start editItem" as default first message
        - If already in bot's chat
            - Prompt for item to edit
            - goto **GetItemToEdit**
    - /start editItem
        - Prompt for item to edit
        - goto **GetItemToEdit**
    - /feedback, /feedback@toGoListBot
        - Prompt for feedback
        - goto **Feedback** 
- *Add Item States*
    - **AddNewSetName**   
    <sup>(expects text message)</sup> 
        - Store item name
        - Prompt for next action
        - goto **ReadyForNextAction**
    - **ReadyForNextAction**    
    <sup>(expects response from reply markup keyboard)</sup>
        - /setAddress
            - Prompt for Address
            - goto **AddNewSetAddress**
        - /setNotes
            - Prompt for Notes
            - goto **AddNewSetNotes**
        - /setURL
            - Prompt for URL
            - goto **AddNewSetURL**
        - /addImage
            - Prompt for image
            - goto **AddNewSetImages**
        - /addTag
            - Send existing tags to add
            - goto **AddNewSetTags**
        - /removeTag
            - Send existing tags available to remove
            - goto **AddNewRemoveTags**
        - /preview
            - Send existing item data
            - Prompt for next action
        - /submit
            - Prompt submission confirmation
            - goto **ConfirmAddItemSubmit**
        - /cancel
            - goto **Idle**
	- **AddNewSetAddress**  
    <sup>(expects text message)</sup> 
        - Store address
        - Prompt for next action
        - goto **ReadyForNextAction**
	- **AddNewSetNotes**  
    <sup>(expects text message)</sup> 
        - Store notes
        - Prompt for next action
        - goto **ReadyForNextAction**
	- **AddNewSetURL**  
    <sup>(expects text message)</sup> 
        - Store URL
        - Prompt for next action
        - goto **ReadyForNextAction**
	- **AddNewSetImages**  
    <sup>(expects image)</sup> 
        - Store image ID
        - Prompt for next action
        - goto **ReadyForNextAction**
	- **AddNewSetTags**  
    <sup>(expects text message or callback from inline keyboard)</sup> 
        - *text message OR selected existing tag*
            - Store tag
        - /done
            - Prompt for next action
            - goto **ReadyForNextAction**
	- **AddNewRemoveTags**  
    <sup>(expects callback from inline keyboard)</sup> 
        - *Selected existing tag*
            - Remove tag
            - Send remaining tags
        - /done
            - Prompt for next action
            - goto **ReadyForNextAction**
	- **ConfirmAddItemSubmit**  
    <sup>(expects callback from inline keyboard)</sup> 
        - yes
            - Store item in chat's list
            - goto **Idle**
        - no
            - Prompt for next action
            - goto **ReadyForNextAction**
- *Delete Item States*
    - **DeleteSelect**  
    <sup>(expects response from inline keyboard)</sup>
        - Get name of item to delete
        - Prompt delete confirmation
        - goto **DeleteConfirm**
	- **DeleteConfirm**  
    <sup>(expects response from inline keyboard)</sup>
        - yes
            - Delete item
        - no
            - Cancel process
        - goto **Idle**
- *Edit Item States*
	- **GetItemToEdit**  
        <sup>(expects response from inline keyboard)</sup>
        - Get name of item to edit
        - Prompt for next action (leverage AddItem Process)
        - goto **ReadyForNextAction**
- *Query States*
	- **QuerySelectType**  
    <sup>(expects callback from inline keyboard)</sup>
        - /getOne
            - Set QueryNum to 1
            - Prompt for search method. (random/tags/name)
            - goto **QueryOneTagOrName**
        - /getFew
            - Prompt for number of items to get
            - goto **QueryFewSetNum**
        - /getAll
            - Set QueryNum to total number of items
            - goto **QuerySetTags**
	- **QueryOneTagOrName**  
    <sup>(expects response from inline keyboard)</sup>
        - /random
            - goto **QueryRetrieve** 
        - /withTag
            - Send existing tags for selection
            - goto **QuerySetTags**
        - /withName
            - Send names of existing items for selection
            - goto **QueryOneSetName**
	- **QueryOneSetName**  
    <sup>(expects response from inline keyboard)</sup>
        - Get name of item to retrieve
        - Prompt if want images
        - goto **QueryRetrieve**
    - **QueryFewSetNum**  
    <sup>(expects a number)</sup>
        - Set QueryNum to input number
        - Send existing tags for selection
        - goto **QuerySetTags**
	- **QuerySetTags**  
    <sup>(expects response from inline keyboard)</sup>
        - *Existing Tag*
            - Add selected tag for query
        - /done
            - Prompt if want images
            - goto **QueryRetrieve**
	- **QueryRetrieve**  
    <sup>(expects response from inline keyboard)</sup>
        - yes
            - Send items with images
        - no
            - Send items without images
        - goto **Idle**
- *Feedback States*
    - **Feedback**  
    <sup>(expects text message)</sup>
        - Get and store feedback
        - goto **Idle**

	
    

`
