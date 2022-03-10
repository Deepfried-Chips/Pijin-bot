local discordia = require('discordia')
local http = require("http")
local fs = require("fs")
local client = discordia.Client()

local prefix = "pigeon "

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
	if processedmessage == prefix .. 'ping' then
		message.channel:send({content = 'Pong! ' .. message.author.username, reference = {message = message, mention = false}})
	end
	if splitstring(processedmessage," ")[1] .. splitstring(processedmessage," ")[2] == prefix .. 'echo' then
		print("detected")
		if message.member.hasPermission(message.mentionedChannels[1],0x00002000) then
			local send = splitstring(message)
			table.remove(send,1)
			local concatsend = table.concat(send," ")
			message.mentionedChannels[1]:send({content = concatsend})
		end
	end
end)

-- client:on('presenceUpdate',function(member)
-- 	local user = member.user
-- 	if user.bot then
-- 		return
-- 	end
-- 	local activity
-- 	if member.activity ~= nil then
-- 		activity = member.activity
-- 	else
-- 		return
-- 	end
-- 	print(activity.name)
-- 	local activityname
-- 	if type(activity.name) == "string" then
-- 		activityname = string.lower(activity.name)
-- 	end

-- 	if activityname == "genshin impact" then
-- 		local time = os.time()
-- 		local playtime = time - activity.start
-- 		if playtime >= 1800 then
-- 			local guild = member.guild
-- 			local id = member.user.id

-- 			guild.kickUser(id,"played genshin for more than 30 minutes")
-- 		end
-- 	end
-- end)

local fd = fs.openSync("token.txt","r",0)
local token = fs.readSync(fd,4096,0)

client:run('Bot ' .. token)

