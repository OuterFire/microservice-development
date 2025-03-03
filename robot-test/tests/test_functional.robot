*** Settings ***

Library    RequestsLibrary
Library    Collections
Library    ./python_library/helper/utils.py
Library    ./python_library/RedisHelper.py

Resource   resource/keyword_helper.resource

*** Test Cases ***

Test Sample
    [Documentation]    sample test
    ${resp} =    Sample Test
    Should Be Equal As Strings
    ...    ${resp}
    ...    Hello World
    Log    sample test log message: ${resp}

Test default api endpoint
    [Documentation]    rest api get default endpoint

    ${resp} =    GET
    ...    url=${EVENT_PRODUCER_URL}/
    ...    expected_status=404
    Status Should Be
    ...    404
    ...    ${resp}
    Should Be Equal As Strings
    ...    ${resp.json()}[title]
    ...    Endpoint Not Found

Test event api endpoint
    [Documentation]    Event Notification (EN) interface, post a valid event notification to the /event endpoint

    Send Event Message
    ...    {"description":"hello world"}
    ...    200

Test invalid request body event api endpoint
    [Documentation]    Event Notification (EN) interface, post a invalid event notification body to the /event endpoint
    Send Invalid Event Message
    ...    {"description":123}
    ...    400

Test invalid request method event api endpoint
    [Documentation]    Event Notification (EN) interface, post a invalid event notification request method to the /event endpoint
    Send Invalid Event Message
    ...    {"description":"123"}
    ...    405
    ...    ${False}

Test read notification message from redis stream
    [Documentation]    Send an Event Notification using the Event Notification (EN) interface which will generate
    ...    a Notification Message. The Notification Message (NM) interface will send the (NM) to a redis stream.
    ...    Verify the Notification Message structure from the redis stream.

    Send Event Message
    ...    {"description":"test_case_running"}
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
        ...    test_case_running
        #Log    ID: ${item}[id], DESCRIPTION: ${item}[description]
    END