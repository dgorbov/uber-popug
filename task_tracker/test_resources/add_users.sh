# requires kafkacat, see: https://github.com/edenhill/kcat
kafkacat -b localhost:9091 -P -t new-user -T -l ./new_user_data