## Logic for query
case constants.QuerySelectType:
    message, _, err := utils.GetMessage(update)
    if err != nil {
        log.Printf("error setting message: %+v", err)
    }
    switch message {
    case "/getOne":
    case "/getFew":
    case "/getAll":
    default:
        sendTemplateReplies(update, "Please select a response from the provided options")
    }
    return
    // getOne markup (/withTag, /withName), GoTo QueryOneTagOrName
    // getFew GoTo QueryFewSetNum. Message how many they want?
    // getAll GoTo QueryAllRetrieve

/* Ask to get one using tag or name */
case constants.QueryOneTagOrName:
    // withTag inline (tags, /done), GoTo QuerySetTags
    // Send message "Don't add any to consider all places"

    // withName inline (names)

/* Ask for tags to search with */
case constants.QueryOneSetTags:
    // tag addTag, preview current, inline (show tags not yet added, /done)
    // done GoTo QueryOneRetrieve. Markup("yes, no"), ask with pics

/* Ask for name to search with */
case constants.QueryOneSetName:
    // set name, GoTo QueryOneRetrieve. Markup("yes, no"), ask with pics

/* Ask whether want pics, and retrieve */
case constants.QueryOneRetrieve:
    // if name != "", get and show place data.
    // if len(tags) == 0, get all, randomly choose one
    // if len(tags) > 0, get all, extract with matching tags. randomly select one
    // if "yes", send pics, goto Idle

/* Ask how many they want */
case constants.QueryFewSetNum:
    // Set number, GoTo QueryFewSetTags. inline (tags, /done)
    // Send message "Don't add any to consider all places"

/* Ask for tags to search with */
case constants.QueryFewSetTags:
    // tag addTag, preview current, inline (show tags not yet added, /done) 
    // done GoTo QueryFewRetrieve, Markup("yes, no") ask with pics

/* Ask whether want pics and retrieve */
case constants.QueryFewRetrieve:
    // if len(tags) == 0, randomly select 
    // if "yes", send pics. GoTo Idle

/* Ask whether want pics and retrieve */
case constants.QueryAllRetrieve:
    // fetch all, send one by one
    // if "yes", send pics. GoTo Idle
}