local discordia = require('discordia')
local magick = require('magick')
local http = require("http")
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
	client:setStatus("In testing")
end)

client:on('messageCreate', function(message)
	if message.author.id == client.user.id then return end
	print(message.content)
	local processedmessage = string.lower(message.content)
	if processedmessage == prefix .. 'ping' then
		message.channel:send({content = 'Pong! ' .. message.author.username, reference = {message = message, mention = false}})
	end

end)

client:run('Bot ODk4MjIxNDEyMjQ5MTM3MTcy.YWhD4A.1zxfZ6I7PqIOcA2rgzkBP1W8DRQ')

