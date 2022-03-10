local discordia = require('discordia')
local http = require("http")
local fs = require("fs")
local client = discordia.Client()
local enums = discordia.enums

local prefix = "pigeon"

function helpcommand(message)
	message:reply("Haha currently there's no commands beside this one")
end


function splitstring(inputstr, sep)
        if sep == nil then
                sep = "%s"
        end
        local t={}
        for str in string.gmatch(inputstr, "([^"..sep.."]+)") do
                table.insert(t, str)
        end
        return t
end

client:on('ready', function()
	print('Logged in as '.. client.user.username)
	client:setGame("In testing")
end)

client:on('messageCreate', function(message)
	if message.author.id == client.user.id then return end
	print(message.content)
	local processedmessage = string.lower(message.content)
	processedmessage = splitstring(processedmessage," ")
	if processedmessage[1] ~= prefix then
		return
	end
	if processedmessage[2] == 'ping' then
		message.channel:send({content = 'Pong! ' .. message.author.username, reference = {message = message, mention = false}})
	end
	if processedmessage[2] == 'echo' then
		if message.member:hasPermission(message.mentionedChannels.first,enums.permission.manageMessages) then
			local send = splitstring(message.content)
		 	print(table.remove(send,1))
		 	print(table.remove(send,1))
			print(table.remove(send,1))
		 	local concatsend = table.concat(send," ")
		 	message.mentionedChannels.first:send({content = concatsend})
		else
			message.member.user:getPrivateChannel():send({content = ":x: \nYou do not have ```Manage Messages``` permission"})
		end
	end
end)

client:on('presenceUpdate',function(member)
	local user = member.user
	if user.bot then
		return
	end
	local activity
	if member.activity ~= nil then
		activity = member.activity
	else
		return
	end
	print(activity.name)
	local activityname
	if type(activity.name) == "string" then
		activityname = string.lower(activity.name)
	end

	if activityname == "genshin impact" then
		local time = os.time()
		local playtime = time - activity.start
		if playtime >= 1800 then
			local guild = member.guild
			local id = member.user.id

			guild.kickUser(id,"played genshin for more than 30 minutes")
		end
	end
end)

local fd = fs.openSync("token.txt","r",0)
local token = fs.readSync(fd,4096,0)

client:run('Bot ' .. token)

