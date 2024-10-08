package io.github.parkwithease.parkeasy.model

open class Event<out T>(private val content: T) {
    var hasBeenHandled = false
        private set

    fun getContentIfNotHandled(): T? {
        return if (hasBeenHandled) {
            null
        } else {
            hasBeenHandled = true
            content
        }
    }

    companion object {
        /**
         * Instantiates an Event that's already been handled. Necessary because StateFlow does not
         * allow null.
         *
         * @param T the content type.
         * @param content the content to initialize with.
         * @return the initialized handled Event.
         */
        fun <T> initial(content: T): Event<T> {
            val event = Event(content)
            event.hasBeenHandled = true
            return event
        }
    }
}
