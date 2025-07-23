userid user_type            username   capabilities provider       avatar_url 
1      hacker|repo|streamer codegirl   1,2,3,4,5    twitch|github  ...
2       ....                nutshadedude 1,2,3      twitch          ...
capability_id capability_name
1             join_hacka
2             

hackathon_id hackathon_name description owner_id
1            something cool 1

hackathon_cap
id capabilities_id hackathon_id
1     

groups
id    user_id_starter
1       1


participants
p_id  u

hackathon_group
id  group_id  hackathon_id

participants
userid, hackathon_id


SELECT h.*, p.* FROM hackathons h, participants p INNER JOIN participants ON participants.hackathon_id WHERE h.id = p.hackathon_tch

# ACTUAL TODOS
- handle refreshing tokens

/hackathon/twitch
/hackathon/twitch/codegirl007
/hackathon/twitch/codegirl007/hackathon-name

/groups/hackathon/create
/groups

SELECT * FROM hackathon WHERE provider = ? AND owner_id = ? INNER JOIN users id ON hackathons.owner_id = users.id;
