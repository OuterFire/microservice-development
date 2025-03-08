*** Settings ***

Library    RequestsLibrary
Library    Collections
Library    ./python_library/RobustnessHelper.py
Library    ./python_library/RedisHelper.py

Resource   robot_resource/keyword_helper.resource


*** Test Cases ***

Test redis-master pod deleted unexpectedly
    [Documentation]    Goal of this test is to delete the redis-master pod, once the pod gets deleted all microservices
    ...    dependent to redis-master pod will lose client connection to the redis-master. It is expected for the redis-master
    ...    pod to create itself automatically and all the dependent microservices to self recover the client connection.
    ...    Finally send an Event Notification using the Event Notification (EN) interface which will generate
    ...    a Notification Message. The Notification Message (NM) interface will send the (NM) to a redis stream.
    ...    Verify the Notification Message structure from the redis stream.

    ${pod} =    Restart Pod
    ...    ${TEST_NAMESPACE}
    ...    ${REDIS_POD_LABEL}

    Sleep    2s
    Wait For Pod Restart
    ...    ${TEST_NAMESPACE}
    ...    ${pod}

    Sleep    500 milliseconds

    Send Event Message
    ...    {"description":"restated_redis_test"}
    ...    200

    Sleep    500 milliseconds

    @{notification_message_list} =    Get Notification Message From Redis
    ...    ${TEST_NAMESPACE}
    ...    ${REDIS_POD_LABEL}
    ...    ${REDIS_SECRET_NAME}
    ...    ${REDIS_STREAM_NAME}

        FOR    ${notification_message}    IN    @{notification_message_list}
            Dictionary Should Contain Key    ${notification_message}  key=id
            Dictionary Should Contain Key    ${notification_message}  key=description
            Should Be Equal As Strings
            ...    ${notification_message}[description]
            ...    restated_redis_test
            #Log    ID: ${notification_message}[id], DESCRIPTION: ${notification_message}[description]
        END

Test event-producer pod deleted unexpectedly
    [Documentation]    Goal of this test is to delete the event-producer pod. It is expected for the event-producer
    ...    pod to create itself automatically and re-connect to the redis client.
    ...    Finally send an Event Notification using the Event Notification (EN) interface which will generate
    ...    a Notification Message. The Notification Message (NM) interface will send the (NM) to a redis stream.
    ...    Verify the Notification Message structure from the redis stream.

    ${pod} =    Restart Pod
    ...    ${TEST_NAMESPACE}
    ...    ${EVENT_PRODUCER_POD_LABEL}

    Sleep    2s
    Wait For Pod Restart
    ...    ${TEST_NAMESPACE}
    ...    ${pod}

    Send Event Message
    ...    {"description":"restated_event_producer_test"}
    ...    200

    Sleep    500 milliseconds

    @{notification_message_list} =    Get Notification Message From Redis
    ...    ${TEST_NAMESPACE}
    ...    ${REDIS_POD_LABEL}
    ...    ${REDIS_SECRET_NAME}
    ...    ${REDIS_STREAM_NAME}

        FOR    ${notification_message}    IN    @{notification_message_list}
            Dictionary Should Contain Key    ${notification_message}  key=id
            Dictionary Should Contain Key    ${notification_message}  key=description
            Should Be Equal As Strings
            ...    ${notification_message}[description]
            ...    restated_event_producer_test
            #Log    ID: ${notification_message}[id], DESCRIPTION: ${notification_message}[description]
        END

# Test verify stream gets created after restart